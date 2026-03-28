#ifndef TSPL_SDK_API_H
#define TSPL_SDK_API_H

#include "sdk_common_api.h"

#define SDK_VERSION     "1,3,1,15"

SDK_API int CALL_STACK TSPL_SelfTest( void* handle );

SDK_API int CALL_STACK TSPL_BitMap( void* handle, int xPos, int yPos, int mode, const TCHAR* fileName ,int iBrightness, int iContrast,int iHtMode);

SDK_API int CALL_STACK TSPL_BitMapStream(void* handle, int x_pos, int y_pos, int mode, char* data, int len, int iBrightness, int iContrast, int iHtMode);

SDK_API int CALL_STACK TSPL_Setup( void* handle, int labelWidth, int labelHeight, int speed, int density, int type, int gap, int offset );

SDK_API int CALL_STACK TSPL_ClearBuffer( void* handle );

SDK_API int CALL_STACK TSPL_Box( void* handle, int xStart, int yStart, int xEnd, int yEnd, int thickness );

SDK_API int CALL_STACK TSPL_BarCode( void* handle, int xPos, int yPos, int codeType, int height, int readable, int rotation, int narrow, int wide, const TCHAR* data );

SDK_API int CALL_STACK TSPL_QrCode( void* handle, int xPos, int yPos, int eccLevel, int width, int mode, int rotation, int model, int mask, const TCHAR* data );

SDK_API int CALL_STACK TSPL_Dmatrix( void* handle, int x_pos, int y_pos, int width, int heigth, int xm, int row, int col, const TCHAR* data );

SDK_API int CALL_STACK TSPL_Text( void* handle, int xPos, int yPos, int font, int rotation, int xMultiplication, int yMultiplication,int alignment, const TCHAR* content );

SDK_API int CALL_STACK TSPL_TextCompatible(void* handle, int xPos, int yPos, int font, int rotation, int xMultiplication, int yMultiplication,const TCHAR* content);

SDK_API int CALL_STACK TSPL_Print( void* handle, int num, int copies );

SDK_API int CALL_STACK TSPL_FormFeed( void* handle );

SDK_API int CALL_STACK TSPL_SetTear( void* handle, int enable );

SDK_API int CALL_STACK TSPL_SetRibbon( void* handle, int ribbon );

SDK_API int CALL_STACK TSPL_Offset( void* handle, int distance );

SDK_API int CALL_STACK TSPL_Direction( void* handle, int direction );

SDK_API int CALL_STACK TSPL_Feed( void* handle, int len );

SDK_API int CALL_STACK TSPL_Home( void* handle );

SDK_API int CALL_STACK TSPL_Learn( void* handle );

SDK_API int CALL_STACK TSPL_GetDllVersion( void* handle );

SDK_API int CALL_STACK TSPL_GetSN( void* handle, char* sn );

SDK_API int CALL_STACK TSPL_GetPrinterStatus(void* handle, int* status);


SDK_API int CALL_STACK TSPL_SetCodePage( void* handle, char* codepage );

SDK_API int CALL_STACK TSPL_PDF417( void* handle, int xPos, int yPos, int width, int height, int rotate, const TCHAR* option, const TCHAR* data );

SDK_API int CALL_STACK TSPL_Block( void* handle, int xPos, int yPos, int width, int height, int font, int rotate, int xMultiplication, int yMultiplication, int space, int alginment, const char* data );

SDK_API int CALL_STACK TSPL_Reverse( void* handle, int xPos, int yPos, int width, int height );

SDK_API int CALL_STACK TSPL_GapDetect(void* handle, int paperLength, int gapLength);

SDK_API int CALL_STACK TSPL_SetCutterEnable(void* handle, int enable);
SDK_API int CALL_STACK TSPL_SetCutter(void* handle, int copies );
SDK_API int CALL_STACK TSPL_Cut(void* handle, int copies );

SDK_API int CALL_STACK TSPL_Bold(void* handle, int bold );

//add by wjz 20210820
SDK_API int CALL_STACK TSPL_Bar(void* handle, int x_pos, int y_pos, int width, int height);
SDK_API int CALL_STACK TSPL_Diagonal(void* handle, int x_pos, int y_pos, int width, int height, int thickness);
/*
 * SDK_API int CALL_STACK TSPL_PrintImageData(void* handle, int x, int y, int mode, int w, int h,  char* data);
 * SDK_API char*  CALL_STACK PrtGetPrinterID(void* handle,char* id);
 */
 
SDK_API int CALL_STACK GetVersion(char* buffer, int len);

SDK_API int CALL_STACK TSPL_GetFirmwareVersion(void* handle, char* buffer, int len);

#endif

