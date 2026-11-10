import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const rule = el.querySelector('path[d="M12 3v18"]');
  animate(rule, { scaleY: [0, 1.25, 1], opacity: [0, 1] }, { duration: τmid });

  const left = el.querySelector('path[d^="m8 8"]');
  animate(left, { translateX: [0, -2, 0.5, 0] }, { duration: τmid, delay: τmin });
  const right = el.querySelector('path[d^="m16 16"]');
  animate(right, { translateX: [0, 2, -0.5, 0] }, { duration: τmid, delay: τmin });
}
