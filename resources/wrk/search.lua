letters = {"А", "Б", "В", "Г", "Д", "Е", "Ё", "Ж", "З", "И", "Й", "К", "Л", "М", "Н", "О", "П", "Р", "С", "Т", "У",
            "Ф", "Х", "Ц", "Ч", "Ш", "Щ", "Ь", "Ы", "Ъ", "Э", "Ю", "Я"}

function escape (str)
    str = string.gsub (str, "([^0-9a-zA-Z !'()*._~-])", -- locale independent
            function (c) return string.format ("%%%02X", string.byte(c)) end)
    str = string.gsub (str, " ", "+")
    return str
end

request = function()
    l1 = letters[math.random(#letters)]
    l2 = letters[math.random(#letters)]
    path = ("/v1/user/search?first_name=%s&last_name=%s"):format(escape(l1),escape(l2))
    return wrk.format("GET", path)
end