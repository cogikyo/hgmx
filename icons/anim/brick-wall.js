import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const wall = el.querySelector("g");
  animate(wall, { scaleY: [0, 1.1, 1] }, { duration: τmid });

  const joints = ["M8 15v6", "M16 15v6", "M12 9v6", "M8 3v6", "M16 3v6"].map((d) => el.querySelector(`path[d="${d}"]`));
  animate(joints, { scaleY: [0, 1] }, { duration: τmin, delay: stagger(0.06, { startDelay: τmin }) });
}
