import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const battery = el.querySelector("rect");
  animate(battery, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const signs = [el.querySelector('path[d="M6 13h4"]'), el.querySelector("g")];
  animate(signs, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: stagger(0.12, { startDelay: τmin }) });
}
