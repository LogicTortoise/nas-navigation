import puppeteer from 'puppeteer';

const delay = ms => new Promise(resolve => setTimeout(resolve, ms));

async function test() {
  console.log('Launching Chrome...');

  const browser = await puppeteer.launch({
    headless: false,
    executablePath: '/Applications/Google Chrome.app/Contents/MacOS/Google Chrome',
    args: ['--no-sandbox', '--disable-setuid-sandbox']
  });

  const page = await browser.newPage();
  await page.setViewport({ width: 1400, height: 900 });

  console.log('Navigating to http://localhost:3000...');
  await page.goto('http://localhost:3000', { waitUntil: 'networkidle2', timeout: 30000 });
  await delay(2000);

  // 更精确的弹窗检测 - 检查 z-50 的 fixed 元素
  const checkModal = async (expectedText) => {
    return await page.evaluate((text) => {
      // 查找所有 z-50 的 fixed 元素
      const allElements = document.querySelectorAll('*');
      for (const el of allElements) {
        const style = window.getComputedStyle(el);
        if (style.position === 'fixed' && style.zIndex === '50') {
          // 检查是否包含 form 或特定文本
          const hasForm = el.querySelector('form');
          const hasInput = el.querySelector('input');
          const containsText = el.textContent.includes(text);
          if (hasForm || hasInput) {
            return {
              visible: true,
              text: el.textContent.substring(0, 300).replace(/\s+/g, ' '),
              containsExpected: containsText
            };
          }
        }
      }
      return { visible: false, text: '', containsExpected: false };
    }, expectedText);
  };

  // 测试结果记录
  const results = { passed: 0, failed: 0 };

  console.log('\n=== 测试1: 点击顶部"添加"按钮 ===');
  try {
    // 找到 header 中包含"添加"文字的按钮
    const clicked = await page.evaluate(() => {
      const header = document.querySelector('header');
      if (!header) return false;
      const buttons = header.querySelectorAll('button');
      for (const btn of buttons) {
        if (btn.textContent.includes('添加')) {
          btn.click();
          return true;
        }
      }
      return false;
    });

    if (clicked) {
      await delay(1000);

      const result = await checkModal('添加网站链接');
      if (result.visible) {
        console.log('✅ 弹窗已显示');
        console.log('   内容预览:', result.text.substring(0, 80) + '...');
        results.passed++;
      } else {
        console.log('❌ 弹窗未显示');
        results.failed++;
      }

      // 点击取消按钮关闭弹窗 (Escape键不工作)
      await page.evaluate(() => {
        const buttons = document.querySelectorAll('button');
        for (const btn of buttons) {
          if (btn.textContent.includes('取消')) {
            btn.click();
            return;
          }
        }
        // 如果没找到取消按钮，点击backdrop
        const backdrop = document.querySelector('.fixed.inset-0.bg-black\\/40');
        if (backdrop) backdrop.click();
      });
      await delay(500);
    } else {
      console.log('❌ 未找到顶部添加按钮');
      results.failed++;
    }
  } catch (e) {
    console.log('❌ 测试失败:', e.message);
    results.failed++;
  }

  console.log('\n=== 测试2: 点击侧边栏"添加分类"按钮 ===');
  try {
    // 找到包含"添加分类"文字或有fa-plus图标的按钮
    const clickResult = await page.evaluate(() => {
      const aside = document.querySelector('aside');
      if (!aside) return { success: false, reason: 'aside not found' };

      // 找到包含"添加分类"的按钮或带有+图标的最后一个按钮
      const buttons = aside.querySelectorAll('button');
      let targetBtn = null;

      for (const btn of buttons) {
        const text = btn.textContent.trim();
        if (text.includes('添加分类')) {
          targetBtn = btn;
          break;
        }
      }

      // 如果没找到带文字的，找最后一个带fa-plus图标的按钮
      if (!targetBtn) {
        for (const btn of buttons) {
          if (btn.querySelector('.fa-plus')) {
            targetBtn = btn;
          }
        }
      }

      if (targetBtn) {
        targetBtn.click();
        return { success: true, btnText: targetBtn.textContent.trim() };
      }
      return { success: false, reason: 'button not found', totalButtons: buttons.length };
    });

    console.log('   按钮点击:', clickResult.success ? '成功' : '失败');
    if (clickResult.btnText) console.log('   按钮文字:', clickResult.btnText);
    if (clickResult.reason) console.log('   原因:', clickResult.reason);

    await delay(1000);

    const result = await checkModal('添加分类');
    if (result.visible && result.text.includes('分类名称')) {
      console.log('✅ 分类弹窗已显示');
      console.log('   内容预览:', result.text.substring(0, 80) + '...');
      results.passed++;
    } else {
      console.log('❌ 弹窗未显示或显示了错误的弹窗');
      console.log('   检测到内容:', result.text?.substring(0, 50));
      await page.screenshot({ path: 'debug-category.png' });
      console.log('   已保存截图: debug-category.png');
      results.failed++;
    }

    // 点击取消按钮关闭弹窗
    await page.evaluate(() => {
      const buttons = document.querySelectorAll('button');
      for (const btn of buttons) {
        if (btn.textContent.includes('取消')) {
          btn.click();
          return;
        }
      }
    });
    await delay(500);
  } catch (e) {
    console.log('❌ 测试失败:', e.message);
    results.failed++;
  }

  console.log('\n=== 测试3: 点击"添加链接"卡片按钮 ===');
  try {
    // 找到包含"添加链接"文字的卡片按钮（不在header中）
    const clicked = await page.evaluate(() => {
      const buttons = document.querySelectorAll('main button');
      for (const btn of buttons) {
        if (btn.textContent.includes('添加链接')) {
          btn.click();
          return true;
        }
      }
      return false;
    });

    console.log('   按钮点击:', clicked ? '成功' : '失败');
    await delay(1000);

    const result = await checkModal('添加网站链接');
    if (result.visible) {
      console.log('✅ 弹窗已显示');
      results.passed++;
    } else {
      console.log('❌ 弹窗未显示');
      results.failed++;
    }

    // 点击取消按钮关闭弹窗
    await page.evaluate(() => {
      const buttons = document.querySelectorAll('button');
      for (const btn of buttons) {
        if (btn.textContent.includes('取消')) {
          btn.click();
          return;
        }
      }
    });
    await delay(500);
  } catch (e) {
    console.log('❌ 测试失败:', e.message);
    results.failed++;
  }

  console.log('\n=== 测试4: 完整流程 - 添加测试分类 ===');
  try {
    // 点击添加分类
    await page.evaluate(() => {
      const aside = document.querySelector('aside');
      const nav = aside?.querySelector('nav');
      const buttons = nav?.querySelectorAll('button') || [];
      buttons[buttons.length - 1]?.click();
    });
    await delay(1000);

    // 输入分类名
    const nameInput = await page.$('input[placeholder*="分类名称"], input[maxlength="20"]');
    if (nameInput) {
      const testCategoryName = 'E2E测试' + Date.now();
      await nameInput.type(testCategoryName);
      console.log('   输入分类名:', testCategoryName);

      // 点击保存
      const saved = await page.evaluate(() => {
        const buttons = document.querySelectorAll('button');
        for (const btn of buttons) {
          if (btn.textContent.includes('保存')) {
            btn.click();
            return true;
          }
        }
        return false;
      });

      if (saved) {
        console.log('   点击保存按钮');
        await delay(2000);

        // 刷新页面验证数据是否持久化
        await page.reload({ waitUntil: 'networkidle2' });
        await delay(1500);

        // 验证分类是否创建
        const categoryExists = await page.evaluate((name) => {
          const aside = document.querySelector('aside');
          return aside?.textContent.includes(name.substring(0, 10)) || false;
        }, testCategoryName);

        if (categoryExists) {
          console.log('✅ 分类创建成功并已持久化');
          results.passed++;
        } else {
          console.log('❌ 分类创建失败 - 刷新后未找到');
          results.failed++;
        }
      }
    } else {
      console.log('❌ 未找到输入框');
      results.failed++;
    }
  } catch (e) {
    console.log('❌ 测试失败:', e.message);
    results.failed++;
  }

  console.log('\n========================================');
  console.log(`测试结果: ${results.passed} 通过, ${results.failed} 失败`);
  console.log('========================================');

  await delay(2000);
  await browser.close();

  // 返回适当的退出码
  process.exit(results.failed > 0 ? 1 : 0);
}

test().catch(e => {
  console.error('测试出错:', e);
  process.exit(1);
});
