# Architecture and Design:

Database:  
    - Created 2 tables, `clients` and `documents`  
    - Each client has many documents, each document has 1 client  
    - Documents table stores typical document data such as size, doc type, and status, as well as a `storage_key` that is the pointer to the phyiscal file which gets stored on disk  
    - Currently, each build will attempt to insert the "default" user for the sake of this assignment. I have included an `ON CONFLICT` clause to prevent duplicate-key errors if this gets built more than once.  

API:  
    - This is a small REST API with 3 endpoints:  
        POST api/documents  
        GET api/documents  
        GET api/client  
    - POST to api/documents stores the uploaded file to disk, and writes it's metadata to a row in the database  
    - GET api/documents retrieves documents to show on both the dashboard and admin views. For the dashboard, it takes a clientId param to scope the docs to one client, and for the admin it is scoped to all documents.  
    - GET api/client allows us to retrieve a client's information by their email address  

Client:  
    - I used react-router's BrowserRouter to create a 1 page app which can toggle between client dashboard and solicitor admin  
    - Dashboard view allows clients to upload documents, and see a table of their uploads
    - Admin view allows Solicitors to view uploads, and their status  
    - Backend calls are routed through `api.ts`, which is proxied to the Go server on `:3000` so we don't have to worry about CORS for this assignment  
    - The "default" user's email is hardcoded, to mock a "logged in" state. This is then used to fetch the client's data from the backend for use in doc upload / display of uploaded docs.  

# Assumptions & Trade-offs:  

All users are assumed to be "logged in" and authorized. There is a toggle on the frontend which allows the user to switch between Dashboard and Admin views, which is purely for ease of viewing both in this project. A production system would include complete seperation between the two, perhaps even with seperate login portals, and role-based permissions on the admin side.  

I assumed that document "status" is either being assessed by a 3rd party, or handled elsewhere, and therefore there is no way to manually "verify" or "reject" a document in this app.  

I used react-router's BrowserRouter for speed of setup, as opposed to file-based routing. This is discussed further in the improvements section.  

# Improvements:  

File based routing with SSR and loader functions would be a much more ideal frontend set up. This would allow us to better organize routes, load data before components, and speed up first paint. This is a must for a scaling app, vs the current BrowserRouter system. Loading client data post auth would then replace the hard coded values in the front and backend of this project.  

Better file validation is needed. I currently am only enforcing a max file size, and doc types via magic byte sniffing. There could be a multitude of checks including strict file size limits per type, visual checks for images to prevent blurry pictures from being accepted, and resolution limits.  

Files should be encrypted in transit and in storage.  

All sensitive docs should be stored / kept for as little time as possible. If we could store them while pending, and then keep a record of their verification once it is complete, we should then reference that record and delete the sensitive docs. This would greatly reduce business liability.  

Storage of files should be done with a cloud provider as opposed to on disk  

Enforce 2FA, email verification for all clients and staff  

The UI is bare bones. Improved styling, component placement, verbiage, error dispaly, etc. would need to happen for a production ready app.  

# Building and Running This Project:  

After cloning this repository, you will need to open 3 terminal windows.

1. `cd` into project root, and run `docker compose -f db/compose.yaml up -d` to build and run the database from it's docker file.  

2. From a second terminal, run `cd server` to move into the server folder. Then run `go mod tidy`, and then `go run .`.  

3. From the third terminal, run `cd client` to move into the client folder. Then run `npm install`, and then `npm run dev`. 

At this point you should be about to open and view the project at `http://localhost:5173`.  