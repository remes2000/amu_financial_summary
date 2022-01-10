const fs = require('fs')
const readline = require('readline');

const args = process.argv.slice(2);
const outputFile = 'output.json';
let insertSeparator = false;

run();

async function run() {
    if(args.length === 0) {
        console.log('Provide input files');
        return;
    }
    
    fs.writeFileSync(outputFile, '[\n');
    for (inputFile of args) {
        console.log(`Processing ${inputFile}...`);
        await processInputFile(inputFile);
    }
    fs.appendFileSync(outputFile, ']');
}

async function processInputFile(inputFile) {
    const inputFileStream = fs.createReadStream(inputFile);
    const rl = readline.createInterface({input: inputFileStream, crlfDelay: 'Infinity'});
    for await (const line of rl) {
        const values = line.split('\t');
        if (values.length !== 9) {
            throw `Line should contain 9 values`;
        }
        printJsonEntry(values)
    }
}

function printJsonEntry(values) {
    const object = {
        date: values[1],
        title: `${values[2]} ${values[3]} ${values[4]}`,
        amount: values[5]
    };
    fs.appendFileSync(outputFile, '\t' + (insertSeparator?',':'') + JSON.stringify(object) + '\n' );
    insertSeparator = true;
}