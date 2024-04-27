import { useRef, useState } from 'react';
import { Grid , Stack, Button } from '@mui/material';
import Editor from '@monaco-editor/react';

import { DrawerComponent } from '../Layout';
import { useCode } from '../hooks';

export const Console = () => {
  const editorRef = useRef(null);
  const [ codeValue, setCodeValue ] = useState('');
  const { output, setActiveCode, setCodeOutput, deleteAll } = useCode();

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

  const handleDeleteAll = () => {
    deleteAll();
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
            height="60vh"
            width="100%"
            theme='vs-dark'
            sx={{ borderRadius: 2 }}
          />
        </Grid>

        <Grid item xs={12}>
          <Stack direction="row" spacing={2}>
            <Editor
              height="30vh"
              defaultValue={output}
              value={output}
              width="90%"
              theme='vs-dark'
              sx={{ borderRadius: 2 }}
              options={{
                readOnly: true
              }}
            />
            <Stack
              direction="column"
              spacing={2}
            >
              <Button
                sx={{
                    backgroundColor: '#ffc300',
                    color: '#000000',
                    '&:hover': {
                        backgroundColor: '#ffc300'
                    },
                    '&:focus': {
                        backgroundColor: '#ffc300'
                    },
                    '&:active': {
                        backgroundColor: '#ffc300'
                    }                        
                }}
                onClick={handleExecuteCode}
              >
                Ejecutar
              </Button>
              <Button
                sx={{
                    backgroundColor: '#ffd60a',
                    color: '#000000',
                    '&:hover': {
                        backgroundColor: '#ffd60a'
                    },
                    '&:focus': {
                        backgroundColor: '#ffd60a'
                    },
                    '&:active': {
                        backgroundColor: '#ffd60a'
                    }                        
                }}
                onClick={handleDeleteAll}
              >
                Eliminar
              </Button>
            </Stack>
          </Stack>
        </Grid>
      </Grid>
    </DrawerComponent>
  )
}
