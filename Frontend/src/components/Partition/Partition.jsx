import { Grid, Stack, Typography, IconButton } from '@mui/material';
import PropTypes from 'prop-types'
import SourceOutlinedIcon from '@mui/icons-material/SourceOutlined';

export const Partition = (props) => {
    const { part } = props;

    const handlePartition = (part) => {
        console.log(`Particion ${part.name} seleccionada`);
        console.log(`Usuarios disponibles ${part.users} seleccionados`)
        console.log(`Sistema de archivos ${part.fileSystem} seleccionado`)
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
