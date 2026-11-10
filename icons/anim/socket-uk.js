import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const plate = el.querySelector("rect");
  animate(plate, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const slots = ["M12 8v2", "M10 15H8", "M14 15h2"].map((d) => el.querySelector(`path[d="${d}"]`));
  animate(slots, { scale: [0, 1.25, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06, { startDelay: τmin }) });
}
