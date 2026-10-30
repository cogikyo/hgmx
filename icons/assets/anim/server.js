import { animate, stagger, τmin } from "../motion.js";

export default function (el) {
  const racks = el.querySelectorAll("rect");
  animate(racks, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: 0.3, delay: stagger(0.09) });

  const lights = el.querySelectorAll("line");
  animate(lights, { scale: [0, 1.75, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.09, { startDelay: τmin }) });
}
