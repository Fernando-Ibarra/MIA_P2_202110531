import { useState } from 'react';
import { styled} from "@mui/material/styles";
import { Box, Toolbar, IconButton, Typography, Drawer, Divider, List, ListItem, ListItemButton, ListItemText, ListItemIcon, Link } from '@mui/material';
import { Link as RouterLink } from 'react-router-dom';
import MuiAppBar from "@mui/material/AppBar";
import CssBaseline from "@mui/material/CssBaseline";
import PropTypes from 'prop-types';

import MenuIcon from '@mui/icons-material/Menu';
import ChevronLeftIcon from '@mui/icons-material/ChevronLeft';
import FolderCopyRoundedIcon from '@mui/icons-material/FolderCopyRounded';
import PictureAsPdfRoundedIcon from '@mui/icons-material/PictureAsPdfRounded';
import CheckBoxOutlineBlankRoundedIcon from '@mui/icons-material/CheckBoxOutlineBlankRounded';

const drawerWidth = 260;

const sites = [
  {
    uuid: '1',
    name: 'Consola',
    icon: CheckBoxOutlineBlankRoundedIcon,
    path: '/',
  },
  {
    uuid: '2',
    name: 'Sistema de Archivos',
    icon: FolderCopyRoundedIcon,
    path: '/file-system',
  },
  {
    uuid: '3',
    name: 'Reportes',
    icon: PictureAsPdfRoundedIcon,
    path: '/reports',
  },
];

const Main = styled("main", { shouldForwardProp: (prop) => prop !== "open" })(
  ({ theme, open }) => ({
    flexGrow: 1,
    padding: theme.spacing(3),
    transition: theme.transitions.create("margin", {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.leavingScreen,
    }),
    marginLeft: `-${drawerWidth}px`,
    ...(open && {
      transition: theme.transitions.create("margin", {
        easing: theme.transitions.easing.easeOut,
        duration: theme.transitions.duration.enteringScreen,
      }),
      marginLeft: 0,
    }),
  })
);

const AppBar = styled(MuiAppBar, {
  shouldForwardProp: (prop) => prop !== "open", })(({ theme, open }) => ({
  transition: theme.transitions.create(["margin", "width"], {
    easing: theme.transitions.easing.sharp,
    duration: theme.transitions.duration.leavingScreen,
  }),
  ...(open && {
    width: `calc(100% - ${drawerWidth}px)`,
    marginLeft: `${drawerWidth}px`,
    transition: theme.transitions.create(["margin", "width"], {
      easing: theme.transitions.easing.easeOut,
      duration: theme.transitions.duration.enteringScreen,
    }),
  }),
}));

const DrawerHeader = styled("div")(({ theme }) => ({
  display: "flex",
  alignItems: "center",
  padding: theme.spacing(0, 1),
  ...theme.mixins.toolbar,
  justifyContent: "flex-end",
}));


export const DrawerComponent = ({ children, name }) => {
    const [open, setOpen] = useState(false);

    const handleDrawerOpen = () => {
      setOpen(true);
    };
    
    const handleDrawerClose = () => {
      setOpen(false);
    };

    return (
        <Box
          sx={{
            display: 'flex',
          }}
        >
          <CssBaseline />
          <AppBar position="fixed" open={open}
            sx={{
              backgroundColor: '#003566',
              height: '80px',
            }}
          >
            <Toolbar>
              <IconButton
                color="inherit"
                aria-label="open drawer"
                onClick={handleDrawerOpen}
                edge="start"
                sx={{ mr: 2, ...(open && { display: "none" }) }}
              >
                <MenuIcon />
              </IconButton>
              <Typography variant="h6" noWrap component="div">
                USAC DRIVE - {name} 
              </Typography>
            </Toolbar>
          </AppBar>
          <Drawer
            sx={{
              width: drawerWidth,
              flexShrink: 0,
              '& .MuiDrawer-paper': {
                width: drawerWidth,
                backgroundColor: '#001d3d',
                boxSizing: 'border-box',
              },
            }}
            variant="persistent"
            anchor="left"
            open={open}
          >
            <DrawerHeader>
              <IconButton onClick={handleDrawerClose}>
                <ChevronLeftIcon 
                  sx={{ color: 'white' }}
                />
              </IconButton>
            </DrawerHeader>
            <Divider
              sx={{
                backgroundColor: 'white'
              }}
            />
            <List>
              { sites.map((site) => (
                <ListItem key={site.uuid} disablePadding>
                  {
                    (site.name == name) ? (
                      <Link
                        to={site.path}
                        component={RouterLink}
                        underline='none'
                        sx={{
                          margin: '0',
                          padding: '0',
                          width: '100%',
                        }}
                      >
                        <ListItemButton
                          sx={{
                            backgroundColor: '#ffc300',
                            '&:hover': {
                              backgroundColor: '#ffd60a'
                            },
                            '&:focus': {
                              backgroundColor: '#ffd60a'
                            }
                          }}
                        >
                          <ListItemIcon>
                                <site.icon
                                  sx={{ color: '#003566' }}
                                />
                              </ListItemIcon>
                              <ListItemText primary={site.name}
                                sx={{ color: '#003566' }}
                              />

                        </ListItemButton>
                      </Link>
                    ) : (
                      <Link
                          to={site.path}
                          component={RouterLink}
                          underline='none'
                          sx={{
                            margin: '0',
                            padding: '0',
                            width: '100%',
                          }}
                        >  
                        <ListItemButton>
                          <ListItemIcon>
                              <site.icon
                                sx={{ color: 'white' }}
                              />
                            </ListItemIcon>
                            <ListItemText primary={site.name}
                              sx={{ color: 'white' }}
                            />
                          
                        </ListItemButton>
                        </Link>
                    )
                  }
                </ListItem>
              ))}
            </List>
          </Drawer>
          <Main
            open={open}
          >
            <DrawerHeader />
            { children }
          </Main>
        </Box>
    )
}


DrawerComponent.propTypes = {
    children: PropTypes.node,
    name: PropTypes.string
};