import { animate, stagger, τmid } from "../motion.js";

export default function (el) {
  const handset = el.querySelector('path[d^="M13.832"]');
  animate(handset, { rotate: [0, -10, 6, 0] }, { duration: τmid });

  const waves = [el.querySelector('path[d^="M13 6"]'), el.querySelector('path[d^="M13 2"]')];
  animate(waves, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.09) });
}
