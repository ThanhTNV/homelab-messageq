# Connect to localhost:8080, send "Hello Server", read response, print, then close
$client = New-Object System.Net.Sockets.TcpClient("localhost", 8080)
$stream = $client.GetStream()
$writer = New-Object System.IO.StreamWriter($stream)
$reader = New-Object System.IO.StreamReader($stream)

# Send a message
$writer.WriteLine("Hello Server")
$writer.Flush()

# Read echoed response
$response = $reader.ReadLine()
Write-Output "Received: $response"

# Close connection
$writer.Close()
$reader.Close()
$client.Close()
