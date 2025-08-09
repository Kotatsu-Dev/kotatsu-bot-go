export const exportJSON = (data: unknown) => {
    let jsonContent = JSON.stringify(data, null, 4);
    let blob = new Blob([jsonContent], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    window.open(url, '_blank');
    URL.revokeObjectURL(url);
}