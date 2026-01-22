module DrawSquare exposing (Params, Bounds, computeBounds)


-- STRUCTURE D’ENTRÉE

type alias Params =
    { size : Int
    , centerLng : Float
    , centerLat : Float
    }


-- STRUCTURE DE SORTIE (bornes Leaflet)

type alias Bounds =
    { southWest : ( Float, Float )
    , northEast : ( Float, Float )
    }


-- COEFFICIENTS
-- 12,5 m ≈ coefficients fournis

lngFactor : Float
lngFactor =
    0.000167

latFactor : Float
latFactor =
    0.000117


-- CALCUL DU CARRÉ

computeBounds : Params -> Bounds
computeBounds p = -- permettre de calculer les 4 coins des carrés 
    { southWest =
        ( p.centerLng - lngFactor / 25 * toFloat p.size
        , p.centerLat - latFactor / 25 * toFloat p.size
        )
    , northEast =
        ( p.centerLng + lngFactor / 25 * toFloat p.size
        , p.centerLat + latFactor / 25 * toFloat p.size
        )
    }
