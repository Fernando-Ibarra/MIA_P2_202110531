import { Grid, Stack, Typography, IconButton } from '@mui/material';
import PropTypes from 'prop-types'
import SourceOutlinedIcon from '@mui/icons-material/SourceOutlined';

export const Partition = (props) => {
    const { part } = props;

    const handlePartition = (part) => {
        console.log(`Particion ${part.name} seleccionada`);
        console.log(`Usuario user: ${part.users[0].user} con pass: ${part.users[0].pass} seleccionado`)
        console.log(`Nombre archivo ${part.fileSystem[1].folder.name} seleccionado`)
        console.log(`Contenido ${part.fileSystem[0].file.content} seleccionado`)
    }

    return (
       <Grid item
            sx={{
                bgcolor: '#003566',
                width: '95%',
                margin: '0',
                borderRadius: '10px',
            }}
            lg={4}
        >
            <Stack
                spacing={2}
                direction="column"
                sx={{
                    justifyContent: 'center',
                    alignItems: 'center',
                }}
            >
                <IconButton
                    onClick={() => handlePartition(part)}
                    sx={{
                        color: 'whites',
                        borderRadius: '10px',
                    }}
                >
                    <SourceOutlinedIcon 
                        sx={{
                            fontSize: '5rem',
                            color: 'white',
                        }}
                    />
                </IconButton>
                <Typography>
                    {part.name}
                </Typography>
            </Stack>
       </Grid>
    )
}

Partition.propTypes = {
    part: PropTypes.object,
}
