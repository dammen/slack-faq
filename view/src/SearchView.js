import React, {useState, useEffect} from 'react';
import './App.css';
import { connect, sendMsg } from "./api/websocket";


function App() {
  const [message, setMessage] = useState("");
  useEffect(() => {
    connect((msg) => {
      console.log(msg)
    });
  })
  return (
    <div className="App">
      <header className="App-header">
        <input placeholder={"input channel name..."} style={stl} value={message} onChange={e => setMessage( e.target.value)}></input>
        <button style={stl} onClick={()=>sendMsg(message)}>Retrieve messages from channel</button>
      </header>
    </div>
  );
}
var stl = {padding: "10px", margin: "10px"}


export default App;
