#ifndef _HPRT_DEFINE_H__
#define _HPRT_DEFINE_H__
#pragma once

#if defined(linux) || defined(__linux) || defined(__linux__) || defined(__TOS_LINUX__)
#	define Platform_Linux
#	define LINUX

#elif defined(_WIN32) || defined(_WIN64)
#	define Platform_Windows
#	define WINDOW

#elif defined(__APPLE__) || defined(TARGET_OS_MAC)
#	define Platform_Mac
#	define MACX

#endif

#ifdef LINUX
#	define DLL_EXPORT	__attribute__ ( (visibility( "default" ) ) )
#	define DLL_EXPORT_WIN	
#else
#	define DLL_EXPORT	__declspec( dllexport )
#	define DLL_EXPORT_WIN	__declspec( dllexport )
#endif

#ifndef HPRT_STATIC_LIB		//DLL
#   ifdef LINUX
#       pragma message("============= soooooooooo")

#	    ifndef SDK_API
#	    	ifdef __cplusplus
#	    		define SDK_API extern"C" DLL_EXPORT
#	    	else
#	    		define SDK_API extern DLL_EXPORT
#	    	endif
#	    endif

#       define CALL_STACK
#       define TCHAR char
#		define COMMON_API_EXPORT

#   else	//WINDOWS

#       pragma message("============= dllllllllll")
#include <tchar.h>
//#       define TCHAR char

#	    ifndef SDK_API
#	    	ifdef __cplusplus
#	    		define SDK_API extern "C"	//DLL_EXPORT  --- 使用模块定义文件以规范__stdcall调用约定下的函数名.
#	    	else
#	    		define SDK_API				//DLL_EXPORT
#	    	endif
#	    endif

#	    ifdef EXPORT_CDECL
#	    	define CALL_STACK __cdecl
#	    	pragma message("============= __cdecl")
#	    else
#	    	define CALL_STACK __stdcall	//windows default
#	    	pragma message("========= __stdcall ")
#	    endif

#   endif   //LINUX

#else		//LIB
#   pragma message("============= static lib")

#	ifndef SDK_API
#		ifdef __cplusplus
#			define SDK_API extern "C" 
#		else
#			define SDK_API 
#		endif
#	endif

#       define CALL_STACK

#   ifndef LINUX
#       include <tchar.h>
#else
#       define TCHAR char
#endif

#endif  //HPRT_STATIC_LIB

#if defined __MINGW32__
#	pragma message("__MINGW32__")
#	ifdef __cplusplus
#		define SDK_API extern "C" __declspec(dllexport)
		#pragma message("__MINGW32__ cxx")
#	else
#		define SDK_API __declspec(dllexport)
#		pragma message("__MINGW32__ c")
#	endif
#endif


#ifdef LINUX
#define TCHAR char
#endif

#ifndef TRUE
#define TRUE 1
#endif
#ifndef FALSE
#define FALSE 0
#endif

#ifndef LINUX
////#define WIN32_LEAN_AND_MEAN
#else
#	define CP_ACP LANG_ENGLISH
#	define Sleep(a) usleep((a) * 1000)
#	define TCHAR char
#endif

#define MAX_PATH    260


#endif	//_HPRT_DEFINE_H__
