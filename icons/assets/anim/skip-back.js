import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const play = el.querySelector('path[d^="M17.971"]');
  animate(play, { translateX: [0, -3, 1, 0] }, { duration: τmax });

  const bar = el.querySelector('path[d="M3 20V4"]');
  animate(bar, { scaleY: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
