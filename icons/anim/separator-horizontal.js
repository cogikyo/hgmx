import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const rule = el.querySelector('path[d="M3 12h18"]');
  animate(rule, { scaleX: [0, 1.25, 1], opacity: [0, 1] }, { duration: τmid });

  const up = el.querySelector('path[d^="m8 8"]');
  animate(up, { translateY: [0, -2, 0.5, 0] }, { duration: τmid, delay: τmin });
  const down = el.querySelector('path[d^="m16 16"]');
  animate(down, { translateY: [0, 2, -0.5, 0] }, { duration: τmid, delay: τmin });
}
