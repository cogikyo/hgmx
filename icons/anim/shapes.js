import { animate, stagger, τmid } from "../motion.js";

export default function (el) {
  const shapes = [el.querySelector("path"), el.querySelector("rect"), el.querySelector("circle")];
  animate(shapes, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid, delay: stagger(0.06) });
}
