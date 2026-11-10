import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const skin = el.querySelector('path[d^="M2.8"]');
  animate(skin, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const rings = el.querySelectorAll("circle");
  animate(rings, { scale: [0, 1.25, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.09, { startDelay: τmin }) });
}
