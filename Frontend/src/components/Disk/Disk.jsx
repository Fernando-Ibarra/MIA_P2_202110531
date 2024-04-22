import { Box, Grid, Stack, Typography } from '@mui/material'
import PropTypes from 'prop-types'

import DiscFullIcon from '@mui/icons-material/DiscFull';
import { Partition } from '../Partition/Partition';

export const Disk = (props) => {
    const { disk } = props;

    return (
        <Box
            sx={{
                bgcolor: '#003566',
                width: '95%',
                margin: '0',
                borderRadius: '10px',
                padding: '10px',
            }}
        >
            <Stack
                sx={{
                    backgroundColor: '#001d3d',
                    bgcolor: '001d3d',
                    justifyContent: 'center',
                }}
            >
                <Stack
                    spacing={2}
                    direction="row"
                    sx={{
                        justifyContent: 'left',
                        alignItems: 'center',
                    }}
                >
                    <DiscFullIcon 
                        sx={{
                            color: '#ffd60a',
                            fontSize: '3rem',
                        }}
                    />

                    <Typography
                        sx={{
                            fontSize: '1.5rem',
                            fontWeight: 'bold',
                            color: 'white',
                        }}
                    >
                        {disk.name}
                    </Typography>
                </Stack>

                <Grid container
                    spacing={2}
                    sx={{
                        justifyContent: 'center',
                        alignItems: 'center',
                        padding: '10px',
                    }}
                    direction="row"
                >
                    {
                        disk.partitions.map((part, index) => (
                            <Partition
                                key={index}
                                part={part}
                            />
                        ))
                    }
                </Grid>
            </Stack>

        </Box>
    )
}

Disk.propTypes = {
    disk: PropTypes.object,
    indexs: PropTypes.number,
}