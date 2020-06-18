<html>
    <head>
        <meta charset="UTF-8">
        <title>goFIDOServer</title>
        <script type="text/javascript" src="http://ajax.googleapis.com/ajax/libs/jquery/2.1.0/jquery.min.js"></script>

        <script type="text/javascript">
            $(function(){
                $("response").html("Response Values");

                $("button").click(function() {
                    var JSONdata = {
                        username: $("#username").val(),
                        displayName: $("#displayName").val()
                    }

                    $.ajax({
                        type: 'post',
                        url: 'http://localhost:8080/attestation/options',
                        data: JSON.stringify(JSONdata),
                        dataType: 'JSON',
                        scriptCharset: 'utf-8',
                        success: function(data) {
                            alert("success");
                            alert(JSON.stringify(data));
                            $("#response").html(JSON.stringify(data));
                        },
                        error: function(data) {
                            alert("error");
                            alert(JSON.stringify(data));
                            $("#response").html(JSON.stringify(data));
                        }
                    });
                })
            })
        </script>
    </head>
    <body>
        <form action="/attestation/options" method="POST">
            <input type="text" id="username" name="username">
            <input type="text" id="displayName" name="displayName">
            <button id="button" type="button">Submit</button>
            <textarea id="response" cols="120" rows="10" disabled></textarea>
        </form>
    </body>
</html>