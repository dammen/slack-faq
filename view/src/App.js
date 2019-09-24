import React from 'react';
import {Link }from "react-router-dom"
import Grid from '@material-ui/core/Grid';
import CssBaseline from '@material-ui/core/CssBaseline';
import Container from '@material-ui/core/Container';
import Typography from '@material-ui/core/Typography';
import Button from '@material-ui/core/Button';
import { makeStyles } from '@material-ui/core/styles';
import Paper from '@material-ui/core/Paper';
const useStyles = makeStyles(theme => ({
  root: {
    flexGrow: 1,
  },
  paper: {
    padding: theme.spacing(2),
    textAlign: 'center',
    color: theme.palette.text.secondary,
  },
}));

function App() {
  const classes = useStyles();

  return (
    <div className={classes.root}>
      <Grid container spacing={3}>
      < Grid item xs={6}>
          <Paper className={classes.paper}>
          <CssBaseline/>
            <Container maxWidth="sm">
              <Typography component="div" style={{ backgroundColor: '#cfe8fc', height: "50vh", alignItems: "center"}}>
                <React.Fragment>
                  <Button><Link to="/list">Go to list view</Link></Button>
                </React.Fragment>
              </Typography>
            </Container>
          </Paper>
        </Grid>
        <Grid item xs={6}>
          <Paper className={classes.paper}>
          <CssBaseline/>
            <Container maxWidth="sm">
              <Typography component="div" style={{ backgroundColor: '#cfe8fc', height: "50vh", alignItems: "center"}}>
                <React.Fragment>
                <Button><Link to="/Search">Go to search view</Link></Button>
                </React.Fragment>
              </Typography>
            </Container>
          </Paper>
        </Grid>
      </Grid>
    </div>
  );
}

export default App;
