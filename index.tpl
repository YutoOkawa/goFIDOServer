<html>
    <head>
        <meta charset="UTF-8">
        <title>goFIDOServer</title>
        <script type="text/javascript" src="https://ajax.googleapis.com/ajax/libs/jquery/2.1.0/jquery.min.js"></script>
        <script type="text/javascript" src="./js/base64url-arraybuffer.js"></script>
        <script type="text/javascript">
            $(function(){
                $("response").html("Response Values");

                $("button").click(function() {
                    var JSONdata = {
                        username: $("#username").val(),
                        displayName: $("#displayName").val()
                    }

                    var attestation = null

                    $.ajax({
                        type: 'post',
                        url: 'https://localhost:8080/attestation/options',
                        data: JSON.stringify(JSONdata),
                        dataType: 'JSON',
                        scriptCharset: 'utf-8',
                        success: function(data) {
                            alert("success");
                            attestation = data;
                            $("#response").html(JSON.stringify(attestation));
                            navigator.credentials.create({ publicKey: attestation })
                            .then(response => {
                                alert("ok");
                            })
                            .catch(error => {
                                alert(error)
                            });
                        },
                        error: function(data) {
                            alert("error");
                            $("#response").html(JSON.stringify(data));
                        }
                    });
                })
            })
        </script>
    </head>
    <body>
        <form action="/attestation/options" method="POST">
            <p><input type="text" id="username" name="username"></p>
            <p><input type="text" id="displayName" name="displayName"></p>
            <p><button id="button" type="button">Submit</button></p>
            <p><textarea id="response" cols="120" rows="10" disabled></textarea></p>
            <p><textarea id="attestation" cols="120" rows="10" disabled></textarea></p>
        </form>
    </body>
</html>