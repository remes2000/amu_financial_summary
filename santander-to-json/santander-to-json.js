const fs = require('fs')
const readline = require('readline');

const args = process.argv.slice(2);
const outputFile = 'output.json';
let insertSeparator = false;

run().then(() => {}).catch(e => console.error(e))

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
    let i = 0;
    for await (const line of rl) {
        i++;
        if(i === 1) {
            continue
        }
        const values = line.split('\t', 8);
        if (values.length !== 8) {
            throw `Line should contain at least 8 values, got ${values.length} ` + line;
        }
        printJsonEntry(values)
    }
}

function printJsonEntry(values) {
    const object = {
        date: values[1],
        title: `${values[2]} ${values[3]} ${values[4]}`,
        amount: values[5].replace(/,/g, '.')
    };
    fs.appendFileSync(outputFile, '\t' + (insertSeparator?',':'') + JSON.stringify(object) + '\n' );
    insertSeparator = true;
}
