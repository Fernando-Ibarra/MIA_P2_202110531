import { Accordion, AccordionDetails, AccordionSummary, Stack, Typography } from '@mui/material';
import PropTypes from 'prop-types'

import ArrowDownwardIcon from '@mui/icons-material/ArrowDownward';

export const Folder = (props) => {
    const { fld } = props;
    return (
        <Stack
            spacing={2}
            direction="column"
            padding={2}
            sx={{
                width: '80%',
            }}
        >
            {
                (Array.isArray(fld))
                ? (
                    fld.map((fol, index) => (
                        (fol.folder )
                        ? (
                            <Accordion
                                key={index}
                                spacing={3}
                            >
                                <AccordionSummary
                                    expandIcon={<ArrowDownwardIcon />}
                                    aria-controls="panel1-content"
                                    id="panel1-header"
                                >
                                    <Typography
                                        sx={{
                                            fontSize: '1.5rem',
                                            fontWeight: 'bold',
                                            color: 'black',
                                        }}
                                    >
                                        {`/${fol.folder.name}`}
                                    </Typography>
                                </AccordionSummary>
                                <AccordionDetails>
                                    {
                                        (fol.folder.content.length >= 1)
                                        ? (
                                            <Folder fld={fol.folder.content} />
                                        )
                                        : null
                                    }
                                </AccordionDetails>
                            </Accordion>
                        )
                        : null
                        
                    ))
                )
                : (
                    <>
                        <Typography
                            sx={{
                                fontWeight: 'regular',
                                color: 'black',
                            }}
                        >
                            {fld}
                        </Typography>
                    </>
                )
            }
        </Stack>
    )
}
Folder.propTypes = {
    fld: PropTypes.array
}
