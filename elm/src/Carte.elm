port module Carte exposing
    ( Model
    , Coord
    , Msg(..)
    , init
    , update
    , view
    , subscriptions
    , initMap
    , drawSquare
    , autoView
    , clearSquares
    , addMarker
    , requestMarker
    )

import Html exposing (Html, div)
import Html.Attributes exposing (id)
import Json.Decode as Decode
import DrawSquare



-- PORTS (anciens + nouveau)

port initMap :
    { lat : Float, lon : Float, zoom : Int }
    -> Cmd msg  -- envoi au JS la commande d'initialisation de la carte à l'aide de lat,lng et zoom


port drawSquare :
    DrawSquare.Bounds
    -> Cmd msg --  -- envoi au JS la commande de traçage des carrés avec en entrée une liste de carrés avec leurs coordonnées


port autoView : () -> Cmd msg -- envoi la commande d'autoview pour que les carrés se voit sur la carte


port clearSquares : () -> Cmd msg -- permet de nettoyer les carrés


port addMarker :
    { lat : Float, lon : Float }
    -> Cmd msg -- envoi un message pour afficher un marqueur à l'endroit qu'on recherche via l'interface


port click_coord : (Decode.Value -> msg) -> Sub msg -- port de reception des coordonnées du point cliqué



-- TYPES

type alias Coord =
    { lat : Float
    , lon : Float
    }


type alias Model =
    { clicked : Maybe Coord }


type Msg
    = Click Decode.Value
    | RequestMarker Coord



-- INIT

init : ( Model, Cmd Msg )
init =
    ( { clicked = Nothing } -- valeur initiale de la carte (centré sur la France)
    , initMap
        { lat = 46.603354
        , lon = 1.888334
        , zoom = 6
        }
    )



-- UPDATE

update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        Click value -> -- si carte cliquée -> on met à jour le model 
            case Decode.decodeValue coordDecoder value of
                Ok coord ->
                    ( { model | clicked = Just coord }, Cmd.none )

                Err _ ->
                    ( model, Cmd.none )

        -- 🔹 demande venant de Main (Interface → Carte) pour récupérer les corrdonnées du point cliqué
        RequestMarker coord ->
            ( { model | clicked = Just coord }
            , addMarker
                { lat = coord.lat
                , lon = coord.lon
                }
            )



-- SUBSCRIPTIONS

subscriptions : Model -> Sub Msg -- -- Écoute les clics sur la carte envoyés par JavaScript
subscriptions _ =
    click_coord Click



-- VIEW

view : Model -> Html msg -- créé la balise HTML de "map"
view _ =
    div [ id "map" ] []



-- DECODER (inchangé)
-- Décode les coordonnées envoyées depuis JavaScript lors d’un clic sur la carte

coordDecoder : Decode.Decoder Coord
coordDecoder =
    Decode.map2 Coord
        (Decode.field "lat" Decode.float)
        (Decode.field "long" Decode.float)



-- API PUBLIQUE POUR MAIN
-- Crée un message destiné à la carte pour demander l’ajout d’un marqueur
requestMarker : Coord -> Msg
requestMarker coord =
    RequestMarker coord
