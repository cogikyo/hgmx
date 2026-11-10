import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const tools = el.querySelectorAll('path[d^="m14.5"], path[d^="M15.634"]');
  animate(tools, { translateY: [0, -2, 0.5, 0] }, { duration: τmid, delay: τmin });

  const tray = el.querySelector('path[d^="M4 12"]');
  animate(tray, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
}
