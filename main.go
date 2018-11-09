package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	pb "github.com/eduborgono/torreControl/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	aero *aeropuerto
}

type avion struct {
	aerolinea string
	vuelo     string
	destino   string
	altura    int32
	pesoMax   int32
	gasMin    int32
}

type aeropuerto struct {
	nombre        string
	pistaAte      [][]*avion
	pistaDes      [][]*avion
	destinos      map[string]string
	avionesNuevos map[string]*avion
	altura        int
	muxDes        sync.Mutex
	muxAte        sync.Mutex
	muxAlt        sync.Mutex
}

func main() {

	insAeropuerto := &aeropuerto{}
	insAeropuerto.configurarAeropuerto()
	insAeropuerto.addDestinos()

	lis, err := net.Listen("tcp", ":7777")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterServicioServer(s, &server{aero: insAeropuerto})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *server) ConsultarDestino(ctx context.Context, in *pb.NuevoAvionRequest) (*pb.NuevoAvionResponse, error) {
	s.aero.avionesNuevos[in.Vuelo] = &avion{aerolinea: in.Linea, vuelo: in.Vuelo, destino: in.Destino, altura: 0, pesoMax: in.Peso, gasMin: in.Combustible}
	s.aero.print("Avión " + in.Vuelo + " quiere despegar")
	s.aero.print("Consultando destino...")
	direccion, exist := s.aero.destinos[strings.ToLower(in.Destino)]
	if exist {
		s.aero.print("Enviando dirección de " + in.Destino + ".")
	} else {
		s.aero.print("No existe ese destino.")
		direccion = "No existe"
	}
	return &pb.NuevoAvionResponse{Direccion: direccion}, nil
}

func (s *server) PedirPermiso(ctx context.Context, in *pb.PermisoRequest) (*pb.PermisoResponse, error) {
	s.aero.print("Consultado restricciones de pasajeros y combustible.")
	var permisoPeso, permisoCombustible bool = false, false
	avion := s.aero.avionesNuevos[in.Vuelo]
	if in.Combustible >= avion.gasMin {
		permisoCombustible = true
	}
	if in.Pasajeros <= avion.pesoMax {
		permisoPeso = true
	}
	return &pb.PermisoResponse{Permiso: permisoPeso && permisoCombustible}, nil
}

func (s *server) PedirInstrucciones(ctx context.Context, in *pb.InstruccionesRequest) (*pb.InstruccionesResponse, error) {
	pistasOcupadas := true
	var pistaSelec int
	var avionPrevio string
	s.aero.muxDes.Lock()
	defer s.aero.muxDes.Unlock()

	for i, pista := range s.aero.pistaDes {
		if len(pista) == 0 {
			pistasOcupadas = false
			pistaSelec = i
			break
		} else {
			menor := s.aero.pistaDes[pistaSelec]
			if len(menor) >= len(pista) {
				pistaSelec = i
				avionPrevio = pista[len(pista)-1].vuelo
			}
		}
	}
	if pistasOcupadas {
		s.aero.print("Todas las pistas ocupadas, enconlando avion.")
	}

	resultado := &pb.InstruccionesResponse{
		PistasOcupadas: pistasOcupadas,
		AvionPrevio:    avionPrevio,
		Pista:          int32(pistaSelec),
		Altura:         int32(s.aero.nuevaAltura()),
	}
	avionAnyadir := s.aero.avionesNuevos[in.Vuelo]
	s.aero.pistaDes[pistaSelec] = append(s.aero.pistaDes[pistaSelec], avionAnyadir)
	delete(s.aero.avionesNuevos, in.Vuelo)
	return resultado, nil
}

func (s *server) VerificarCola(ctx context.Context, in *pb.ColaRequest) (*pb.ColaResponse, error) {
	var nombre string
	if len(s.aero.pistaDes[in.Pista]) == 0 {
		nombre = "Nadie"
	} else {
		nombre = s.aero.pistaDes[in.Pista][0].vuelo
	}
	return &pb.ColaResponse{CabezaCola: nombre}, nil
}

func (s *server) AvisarDespegue(ctx context.Context, in *pb.DespegueRequest) (*pb.DespegueResponse, error) {
	s.aero.muxDes.Lock()
	defer s.aero.muxDes.Unlock()
	if len(s.aero.pistaDes[in.Pista]) > 0 {
		if s.aero.pistaDes[in.Pista][0].vuelo == in.Vuelo {
			dropped := s.aero.pistaDes[in.Pista][0]
			s.aero.pistaDes[in.Pista] = s.aero.pistaDes[in.Pista][1:]
			s.aero.print("Despegó el vuelo " + dropped.vuelo + " de la pista " + strconv.Itoa(int(in.Pista)))
		}
	}
	return &pb.DespegueResponse{}, nil
}

func (s *server) Atterizar(ctx context.Context, in *pb.AterrizajeRequest) (*pb.At_InstruccionesResponse, error) {
	pistasOcupadas := true
	var pistaSelec int
	var avionPrevio string
	s.aero.muxAte.Lock()
	defer s.aero.muxAte.Unlock()

	for i, pista := range s.aero.pistaAte {
		fmt.Println(pista)
		if len(pista) == 0 {
			pistasOcupadas = false
			pistaSelec = i
			break
		} else {
			menor := s.aero.pistaAte[pistaSelec]
			if len(menor) >= len(pista) {
				pistaSelec = i
				avionPrevio = pista[len(pista)-1].vuelo
			}
		}
	}
	if pistasOcupadas {
		s.aero.print("Todas las pistas ocupadas, enconlando avion.")
	}

	alturaAvion := s.aero.nuevaAltura()
	resultado := &pb.At_InstruccionesResponse{
		PistasOcupadas: pistasOcupadas,
		AvionPrevio:    avionPrevio,
		Pista:          int32(pistaSelec),
		Altura:         int32(alturaAvion),
	}
	avionAnyadir := &avion{
		aerolinea: in.Linea,
		vuelo:     in.Vuelo,
		destino:   in.Origen,
		altura:    int32(alturaAvion),
	}
	s.aero.pistaAte[pistaSelec] = append(s.aero.pistaAte[pistaSelec], avionAnyadir)
	return resultado, nil
}

func (s *server) VerificarCola_At(ctx context.Context, in *pb.At_ColaRequest) (*pb.At_ColaResponse, error) {
	var nombre string
	if len(s.aero.pistaAte[in.Pista]) == 0 {
		nombre = "Nadie"
	} else {
		nombre = s.aero.pistaAte[in.Pista][0].vuelo
	}
	return &pb.At_ColaResponse{CabezaCola: nombre}, nil
}

func (s *server) AvisarAterrizaje(ctx context.Context, in *pb.AterrizarRequest) (*pb.AterrizarResponse, error) {
	s.aero.muxAte.Lock()
	defer s.aero.muxAte.Unlock()
	if len(s.aero.pistaAte[in.Pista]) > 0 {
		if s.aero.pistaAte[in.Pista][0].vuelo == in.Vuelo {
			dropped := s.aero.pistaAte[in.Pista][0]
			s.aero.pistaAte[in.Pista] = s.aero.pistaAte[in.Pista][1:]
			s.aero.print("Aterrizó el vuelo " + dropped.vuelo + " en la pista " + strconv.Itoa(int(in.Pista)))
		}
	}
	return &pb.AterrizarResponse{}, nil
}

func horaMinSec() string {
	t := time.Now()
	return strconv.Itoa(t.Hour()) + ":" + strconv.Itoa(t.Minute()) + ":" + strconv.Itoa(t.Second())
}

func (instancia *aeropuerto) nuevaAltura() int {
	instancia.muxAlt.Lock()
	defer instancia.muxAlt.Unlock()
	instancia.altura++
	return instancia.altura
}

func (instancia *aeropuerto) print(str string) {
	fmt.Println("[Torre de control - " + instancia.nombre + " " + horaMinSec() + "] " + str)
}

func (instancia *aeropuerto) configurarAeropuerto() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Bienvenido a la Torre de control")

	fmt.Println("[Torre de control] Nombre del Aeropuerto:")
	text, _ := reader.ReadString('\n')
	instancia.nombre = strings.Replace(text, "\r\n", "", -1)

	fmt.Println("[Torre de control - " + instancia.nombre + "] Cantidad de pistas de aterrizaje:")
	text, _ = reader.ReadString('\n')
	cantPistasAte, _ := strconv.Atoi(strings.Replace(text, "\r\n", "", -1))
	instancia.pistaAte = make([][]*avion, cantPistasAte)
	for i := 0; i < cantPistasAte; i++ {
		instancia.pistaAte[i] = make([]*avion, 0)
	}

	fmt.Println("[Torre de control - " + instancia.nombre + "] Cantidad de pistas de despegue:")
	text, _ = reader.ReadString('\n')
	cantPistasDes, _ := strconv.Atoi(strings.Replace(text, "\r\n", "", -1))
	instancia.pistaDes = make([][]*avion, cantPistasDes)
	for i := 0; i < cantPistasDes; i++ {
		instancia.pistaDes[i] = make([]*avion, 0)
	}

	instancia.avionesNuevos = make(map[string]*avion)
	instancia.altura = 0
}

func (instancia *aeropuerto) addDestinos() {
	re, _ := regexp.Compile(`(.+)\s((?:\d+\.)+\d+)`)
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("[Torre de control - " + instancia.nombre + "] Para agregar destino presione enter")
		fmt.Scanln()
		if instancia.destinos == nil {
			instancia.destinos = make(map[string]string)
		}

		fmt.Println("[Torre de control - " + instancia.nombre + "]  Ingrese nombre y direccion IP del destino (o \"fin\"):")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\r\n", "", -1)
		if text != "fin" {
			result := re.FindStringSubmatch(text)
			if len(result) == 3 {
				instancia.destinos[strings.ToLower(result[1])] = result[2]
			}
			fmt.Println("Destino agregado con éxito")
		} else {
			break
		}
	}
}
