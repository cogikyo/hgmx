import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const bar = el.querySelector('path[d="M21 20V4"]');
  animate(bar, { translateX: [0, -2, 0.5, 0] }, { duration: τmid });

  const play = el.querySelector('path[d^="M13.971"]');
  animate(play, { translateX: [0, -2, 0.5, 0] }, { duration: τmid, delay: τmin });
}
