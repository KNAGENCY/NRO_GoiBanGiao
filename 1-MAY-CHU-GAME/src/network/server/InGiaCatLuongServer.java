package network.server;

public interface InGiaCatLuongServer extends Runnable {
  InGiaCatLuongServer init();

  InGiaCatLuongServer start(int paramInt) throws Exception;

  InGiaCatLuongServer setAcceptHandler(ISessionAcceptHandler paramISessionAcceptHandler);

  InGiaCatLuongServer close();

  InGiaCatLuongServer dispose();

  InGiaCatLuongServer randomKey(boolean paramBoolean);

  InGiaCatLuongServer setDoSomeThingWhenClose(IServerClose paramIServerClose);

  InGiaCatLuongServer setTypeSessioClone(Class paramClass) throws Exception;

  ISessionAcceptHandler getAcceptHandler() throws Exception;

  boolean isRandomKey();

  void stopConnect();
}

/*
 * Location:
 * C:\Users\VoHoangKiet\Downloads\TEA_V5\lib\GirlkunNetwork.jar!\com\girlkun\
 * network\server\InGiaCatLuongServer.class
 * Java compiler version: 8 (52.0)
 * JD-Core Version: 1.1.3
 */