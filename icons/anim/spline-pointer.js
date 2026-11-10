import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const curve = el.querySelector('path[d^="M5 17"]');
  animate(curve, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const pointer = el.querySelector('path[d^="M12.034"]');
  animate(pointer, { translateX: [0, -2, 1, 0], translateY: [0, -2, 1, 0] }, { duration: τmid, delay: τmin });
}
