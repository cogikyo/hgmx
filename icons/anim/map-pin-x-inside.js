import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const pin = el.querySelector('path[d^="M20 10c"]');
  animate(pin, { translateY: [-6, 1, 0], scaleY: [1.15, 0.9, 1] }, { duration: τmid });

  const bars = [el.querySelector('path[d^="m14.5 7.5"]'), el.querySelector('path[d^="m9.5 7.5"]')];
  animate(bars, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: stagger(0.06, { startDelay: τmin }) });
}
