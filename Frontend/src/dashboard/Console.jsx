import { useEffect, useRef } from 'react';
import { Grid , Stack } from '@mui/material';
import { DrawerComponent } from '../Layout';
import Editor, { useMonaco } from '@monaco-editor/react';

function setEditorTheme(monaco) {
  monaco.editor.defineTheme('onedark', {
    base: 'vs-dark',
    inherit: true,
    rules: [
      {
        token: 'comment',
        foreground: '#5c6370',
        fontStyle: 'italic'
      },
      { token: 'constant', foreground: '#e06c75' }
    ],
    colors: {
      'editor.background': '#21252b'
    }
  });
}


export const Console = () => {
  const editorRef = useRef(null);
  const handleEditorDidMount = (editor) => {
    editorRef.current = editor;
  }

  const getEditorValue = () => {
    return editorRef.current.getValue();
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
            defaultLanguage="javascript"
            defaultValue="// Ingrese su código aquí"
            onMount={handleEditorDidMount}
            height="65vh"
            width="100%"
            theme='vs-dark'
            sx={{ borderRadius: 2 }}
            beforeMount={setEditorTheme}
          />
        </Grid>

        <Grid item xs={12}>
          <Stack direction="row" spacing={2}>
            <Editor
              defaultLanguage="javascript"
              defaultValue="// Ingrese su código aquí"
              onMount={handleEditorDidMount}
              height="25vh"
              width="80%"
              theme='vs-dark'
              sx={{ borderRadius: 2 }}
              disabled
            />
            <Stack
              direction="column"
              spacing={2}
            >
              <button onClick={() => console.log(getEditorValue())}>Ejecutar</button>
              <button onClick={() => editorRef.current.setValue('')}>Limpiar</button>
            </Stack>
          </Stack>
        </Grid>
      </Grid>
    </DrawerComponent>
  )
}
