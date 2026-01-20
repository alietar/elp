module UserApi exposing (StartPointData, ServerResponse, Square, fetchSquares)

import Http
import Json.Encode as Encode
import Json.Decode as Decode

type alias Square =
    { size : Int
    , centerLat : Float
    , centerLng : Float
    }

type alias ServerResponse = List Square

type alias StartPointData =
    { lat : Float
    , lng : Float
    , deniv : Float
    }

encodeUser : StartPointData -> Encode.Value
encodeUser startPoint =
    Encode.object [
        ( "lat", Encode.float startPoint.lat ),
        ( "lng", Encode.float startPoint.lng ),
        ( "deniv", Encode.float startPoint.deniv )
    ]

squareDecoder : Decode.Decoder Square
squareDecoder = 
    Decode.map3 Square
        (Decode.field "Size" Decode.int)
        (Decode.field "CenterLng" Decode.float)
        (Decode.field "CenterLat" Decode.float)

responseDecoder : Decode.Decoder ServerResponse 
responseDecoder =
    Decode.list squareDecoder

fetchSquares : StartPointData -> (Result Http.Error ServerResponse -> msg) -> Cmd msg
fetchSquares data toMsg =
    Http.post
        { url = "http://localhost:8080/points" -- Ton URL
        , body = Http.jsonBody (encodeUser data)    -- Ou jsonBody si tu dois envoyer des filtres
        , expect = Http.expectJson toMsg responseDecoder
        }