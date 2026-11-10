import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const handset = el.querySelector('path[d^="M13.832"]');
  animate(handset, { rotate: [0, -10, 6, 0] }, { duration: τmid });

  const entries = [el.querySelector('path[d^="M14.396"]'), el.querySelector('path[d^="M17.343"]')];
  animate(entries, { translateX: [-3, 0], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.12, { startDelay: τmin }) });
}
