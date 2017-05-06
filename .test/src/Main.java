class Main {
    public static void main(String[] args) {
        byte[] b = { 
            (byte) 0xff, 
            (byte) 0x10,
            (byte) 0xaa,
            (byte) 0xb4
        };

        int result = 0;
        for(int i = 0; i < 4; i++)
              result = ((result << 8) | b[i]);

        System.out.println(result);
    }
}
