#ifndef _SDK_COMMON_API_H
#define _SDK_COMMON_API_H

#include "hprt_define.h"

//add by wjz 20240104 并发控制多台打印机时，需要手动初始化一下
SDK_API void CALL_STACK SDKInit();
SDK_API void CALL_STACK SDKDeInit();

SDK_API void CALL_STACK FormatError( int errorNo, int langid, unsigned char* buf, int pos, int bufSize );

SDK_API int CALL_STACK PrinterCreator( void** phandle, const TCHAR* model );

SDK_API int CALL_STACK SetLog(int enable, const TCHAR* path);

SDK_API void* CALL_STACK PrinterCreatorS(const TCHAR* model );

SDK_API int CALL_STACK PrinterDestroy( void* handle );

//==================================

/*
//windows上：[见hprt_io_windows.c]
//	串口："COM1,115200,n,8,1,n"		依次是：串口号、波特率，校验位，数据位，停止位，握手/流控
//  USB自动匹配： "USB"									//支持自动匹配我们的打印机
//	USB："USB,设备名【如HT300】"						//连接指定的打印机
//	USB："USB,USB027"									//连接指定的打印机(端口号)
//	USB："USB,Port_#0008.Hub_#0001"						//连接指定的打印机(位置)
//	USB："USB,vid=0x20d1,pid=0x7007"                    //连接指定的打印机(vid/pid)		新增支持的格式 by wjz 20230216
//	USB："USB,设备sn号"									//连接指定的打印机(sn)			新增支持的格式 by wjz 20230925
//	USB："USB,总线关系"									//连接指定的打印机				新增支持的格式 by wjz 20230925
//	NET: "NET,127.0.0.1,9100"
//	BT:	"BT,127.0.0.1,1234"
//  LPT: "LPT1"
//	FILE: "FILE,文件路径"
//  BUFFER: "BUF"
//================================
//linux上：[见hprt_io_linux.c] 【JSSDK上ANY最好改成对应机型名】
//	串口："COM=ttyS1,115200"
//	USB自动匹配： "USB"                                 //支持自动匹配我们的打印机		新增支持的格式 by wjz 20230925
//	USB："USB,设备名【如USBPrinter】"					//连接指定的打印机
//	USB："USB,ANY,vid=0x18d1,pid=0x010b"				//连接指定的打印机(vid/pid)	
//	USB："USB,ANY,bus=003,addr=24"						//连接指定的打印机(端口号)
//	USB："USB,ANY,设备sn号"								//连接指定的打印机(sn)			新增支持的格式 by wjz 20230925
//	NET: "NET,IP=127.0.0.1,PORT=9100"
//	LPT: "LPT1"
//	FILE: "FILE,文件路径"
//==================================
linux上课设置的参数比较少。形式也不尽相同，全部都是字符串形式，USB描述符如上linux都会比windows多一个设备名
windows上，串口可设置参数有：
	波特率：2400,4800,9600,19200,38400,57600,115200  
	校验位：n,o,e,m.s		分别对应no,odd,even,mark,space
	数据位：4,5,6,7,8
	停止位："1","1.5","2"
	流控："p","c","d","x","n"	详见“hprt_io_serial_windows.c”

linux上只支持设置波特率同上

*/

//================================
SDK_API int CALL_STACK PortOpen( void* handle, const TCHAR* ioSettings );

SDK_API int CALL_STACK DriverPortOpen( void* handle, const TCHAR* driverName );

SDK_API int CALL_STACK PortClose( void* handle );

SDK_API int CALL_STACK DirectIO( void* handle, unsigned char* writeData, unsigned int writeNum, unsigned char* readData, unsigned int readNum, unsigned int* preadedNum );

SDK_API int CALL_STACK WriteData( void* handle, unsigned char* writeData, unsigned int writeNum );

SDK_API int CALL_STACK ReadData( void* handle, unsigned char* readData, unsigned int readNum, unsigned int* preadedNum );

#if defined(_WIN32) || defined(_WIN64) 
SDK_API int CALL_STACK ReadDataClean(void* handle, unsigned char* readData, unsigned int readNum, unsigned int* preadedNum,unsigned int constanttimeout);
#endif

SDK_API int CALL_STACK ReadDataEOF( void* handle,
                                    unsigned char* readData, unsigned int offSet, unsigned int bufLength,
                                    unsigned char soh, unsigned char eof,
                                    unsigned int* preadedNum);

//add by wjz 20231023
SDK_API int CALL_STACK ReadDataEOFEx(void* handle,
									unsigned char* readData, unsigned int offSet, unsigned int bufLength,
									unsigned char soh, unsigned char eof,
									unsigned int* preadedNum,
									int trytimes,
                                    long long int delayms);

SDK_API int CALL_STACK SendCommand( void* handle, char* writeData);

//add by wjz 20220517
SDK_API int CALL_STACK ReadDataExist(void* handle,
									unsigned char* readData, unsigned int offSet, unsigned int bufLength,
									unsigned char* soh, unsigned char* eof,
									unsigned int* preadedNum,
									int trytimes);

//add by wjz 20231023
SDK_API int CALL_STACK ReadDataExistEx(void* handle,
									unsigned char* readData, unsigned int offSet, unsigned int bufLength,
									unsigned char* soh, unsigned char* eof,
									unsigned int* preadedNum,
									int trytimes,
                                    long long int  delayms);


#if defined(_WIN32) || defined(_WIN64)
SDK_API int CALL_STACK GetVidPid(unsigned char* buf);
SDK_API int CALL_STACK GetUsbList(const TCHAR* cUsbList, int buf_size, int* pcnt);
#endif

//add by wjz 20230216 获取设备连接信息 
SDK_API int CALL_STACK GetUriOption(void* handle, char* key, char* value, int buflen);

SDK_API int CALL_STACK SetConfigDir(char* dir, int len);

//add by wjz 20240516
SDK_API int CALL_STACK CheckHandle(void* handle);

//add by wjz 20241024
SDK_API void* CALL_STACK TryLockPrinterTask(const char* model);
SDK_API void* CALL_STACK LockPrinterTask(const char* model);
SDK_API int CALL_STACK UnLockPrinterTask(const char* model, void* m);

#endif
