import { useRef, useState } from 'react';
import { Grid , Stack, Button } from '@mui/material';
import Editor from '@monaco-editor/react';

import { DrawerComponent } from '../Layout';
import { useCode } from '../hooks';

export const Console = () => {
  const editorRef = useRef(null);
  const [ codeValue, setCodeValue ] = useState('');
  const { output, setActiveCode, setCodeOutput } = useCode();

  const handleCodeChange = () => {
    setCodeValue(editorRef.current.getValue());
    setActiveCode(editorRef.current.getValue());
  }
  
  const handleEditorDidMount = (editor) => {
    editorRef.current = editor;
  }

  const handleExecuteCode = () => {
    setCodeOutput();
  }

  return (
    <DrawerComponent
      name="Consola"
    >
      <Grid container
        spacing={2}
        direction="column"
      >
        <Grid item xs={12}>
          <Editor
            defaultValue={codeValue}
            value={codeValue}
            onMount={handleEditorDidMount}
            onChange={handleCodeChange}
            height="65vh"
            width="100%"
            theme='vs-dark'
            sx={{ borderRadius: 2 }}
          />
        </Grid>

        <Grid item xs={12}>
          <Stack direction="row" spacing={2}>
            <Editor
              height="25vh"
              defaultValue={output}
              value={output}
              width="80%"
              theme='vs-dark'
              sx={{ borderRadius: 2 }}
            />
            <Stack
              direction="column"
              spacing={2}
            >
              <Button
                sx={{
                    backgroundColor: '#e87d3e',
                    color: '#000000',
                    '&:hover': {
                        backgroundColor: '#e87d3e'
                    },
                    '&:focus': {
                        backgroundColor: '#e87d3e'
                    },
                    '&:active': {
                        backgroundColor: '#e87d3e'
                    }                        
                }}
                onClick={handleExecuteCode}
              >
                Código
              </Button>
            </Stack>
          </Stack>
        </Grid>
      </Grid>
    </DrawerComponent>
  )
}
