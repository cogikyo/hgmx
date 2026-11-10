import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const plug = el.querySelectorAll('path:not([d^="m18 3"])');
  animate(plug, { translateX: [0, 1.5, -0.5, 0], translateY: [0, -1.5, 0.5, 0] }, { duration: τmid });

  const bolt = el.querySelector('path[d^="m18 3"]');
  animate(bolt, { translateX: [2, -0.5, 0], translateY: [-2, 0.5, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
