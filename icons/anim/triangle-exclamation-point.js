import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const sign = el.querySelector("g");
  animate(sign, { rotate: [0, -10, 8, -4, 0] }, { duration: τmax });
  const bar = el.querySelector('path[d="M12 9v4"]');
  animate(bar, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });
  const dot = el.querySelector('path[d^="M12 17"]');
  animate(dot, { scale: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
