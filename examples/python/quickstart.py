from httpresponder import HTTPResponder

client = HTTPResponder()
print(client.decode('SGVsbG8sIHdvcmxkIQ==').text())
print('Requested status:', client.status(418).status)
