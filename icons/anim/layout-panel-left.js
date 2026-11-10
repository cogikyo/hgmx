import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const panel = el.querySelector('rect[height="18"]');
  animate(panel, { scaleY: [0, 1.1, 1], opacity: [0, 1] }, { duration: τmid });

  const tiles = el.querySelectorAll('rect[height="7"]');
  animate(tiles, { scale: [0, 1.25, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.12, { startDelay: τmin }) });
}
