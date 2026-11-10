import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const panel = el.querySelector('rect[width="18"]');
  animate(panel, { scaleX: [0, 1.1, 1], opacity: [0, 1] }, { duration: τmid });

  const tiles = el.querySelectorAll('rect[width="7"]');
  animate(tiles, { scale: [0, 1.25, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.12, { startDelay: τmin }) });
}
