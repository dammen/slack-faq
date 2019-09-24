import React, {useState} from 'react';
import CssBaseline from '@material-ui/core/CssBaseline';
import Container from '@material-ui/core/Container';
import Typography from '@material-ui/core/Typography';

import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';

function App() {

  const [token, setToken] = useState("");
  const [channels, setChannels] = useState([])

  function fetchChannels(){
    let url = "https://slack.com/api/conversations.list?token=" + token + "&limit=" + 50
    fetch(url)
        .then(res => res.json())
        .then((res) => {
            let channelNames = res.channels.map(channel => channel.name)
            setChannels(channelNames)
        })
    .catch(console.log)
}
  return (
    <React.Fragment>
        <CssBaseline />
        <Container maxWidth="sm">
        <Typography component="div" style={{ backgroundColor: '#cfe8fc' }}>
            <React.Fragment>

            <input placeholder={"input channel token..."} style={stl} value={token} onChange={e => setToken(e.target.value)}></input>
            <button style={stl} onClick={()=>fetchChannels()}>Retrieve messages from channel</button>
            </React.Fragment>

        </Typography>

        </Container>
        <Container maxWidth="sm">
            <Typography component="div" style={{ backgroundColor: '#cfe8fc'}}>
                    <List component="nav" aria-label="Channels">
                        {channels.map((name, index)=> (
                        <ListItem button key={index}>
                            <ListItemText primary={name} />
                        </ListItem>
                        ))}
                    </List>
            </Typography>
        </Container>
    </React.Fragment>
    );
}
var stl = {padding: "10px", margin: "10px"}


export default App;
