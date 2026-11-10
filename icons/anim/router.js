import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const waves = [el.querySelector('path[d^="M17.84"]'), el.querySelector('path[d^="M20.66"]')];
  animate(waves, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.09) });

  const leds = el.querySelectorAll('path[d^="M6.01"], path[d^="M10.01"]');
  animate(leds, { scale: [0, 1.5, 1] }, { duration: τmin, delay: τmid });
}
