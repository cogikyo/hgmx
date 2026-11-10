import { animate, stagger } from "../motion.js";

export default function (el) {
  const rays = ["M2 12h", "m4.929 ", "M12 2v", "M18.5 ", "M21 12h"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(rays, { scale: [0, 1.75, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06) });
}
