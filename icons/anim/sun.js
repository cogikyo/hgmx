import { animate, stagger, τmax } from "../motion.js";

export default function (el) {
  const core = el.querySelector("circle");
  animate(core, { scale: [0.5, 1.25, 1] }, { duration: τmax });

  const rays = el.querySelectorAll("path");
  animate(rays, { scale: [0, 1.75, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.03) });
}
