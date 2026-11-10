import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const oven = el.querySelector('rect[width="20"]');
  animate(oven, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const pane = el.querySelector('rect[width="8"]');
  animate(pane, { scale: [1, 0.85, 1.1, 1] }, { duration: τmid, delay: τmin });
}
