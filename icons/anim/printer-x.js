import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const sheet = el.querySelector('path[d^="M12.531 22"]');
  animate(sheet, { translateY: [-5, 1, 0], opacity: [0, 1] }, { duration: τmid });

  const bars = el.querySelectorAll('path[d^="m16.5"]');
  animate(bars, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: stagger(0.06, { startDelay: τmin }) });
}
