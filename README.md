# CameraFileCopy

Este repositório contém o código-fonte de um programa simples que copia arquivos de uma pasta de origem para uma pasta de destino. Criado como forma de aplicação prática dos meus estudos, com uma finalidade específica, copiar arquivos de câmera para um HD externo, removendo arquivos antigos para liberar espaço.

## Funcionalidades

- Copia arquivos de uma pasta de origem para uma pasta de destino.
- Remove arquivos antigos de uma pasta de origem.
- Mantém metadados, permissões e hierarquia dos arquivos.

## Uso

Para usar o programa, execute o seguinte comando:

```bash
./cameraFileCopy -d <pasta_de_destino> -o <pasta_de_origem> -clean -days <dias> -max <max_itens>
```

Onde:

- <pasta_de_destino>: pasta de destino para onde os arquivos serão copiados.
- <pasta_de_origem>: pasta de origem de onde os arquivos serão copiados.
- \<clean>: opção para limpar arquivos antigos. [opcional]
- \<days>: número de dias para considerar arquivos antigos. [opcional]
- <max_itens>: número máximo de arquivos a serem copiados utilizando Goroutines. [opcional]
