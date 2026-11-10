import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const envelope = el.querySelector('path[d^="M22 12.532"]');
  animate(envelope, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const flap = el.querySelector('path[d^="m22 7"]');
  animate(flap, { scaleY: [-1, 1.15, 1] }, { duration: τmid });

  const bars = [el.querySelector('path[d^="m21.5"]'), el.querySelector('path[d^="m16.5"]')];
  animate(bars, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: stagger(τmin) });
}
